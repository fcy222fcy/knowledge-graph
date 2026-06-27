const puppeteer = require('puppeteer');
const fs = require('fs');
const path = require('path');

// 定义所有功能流程图
const flows = [
    {
        id: '5.1.1-注册功能',
        title: '注册功能实现流程',
        mermaid: `flowchart TD
    A([开始注册]) --> B[填写注册信息<br/>用户名/邮箱/密码/确认密码]
    B --> C{前端表单验证}
    C -->|验证失败| D[显示错误提示]
    D --> B
    C -->|验证通过| E[发送POST请求<br/>/api/v1/auth/register]
    E --> F[Controller接收请求<br/>绑定注册DTO]
    F --> G{参数验证}
    G -->|验证失败| H[返回错误响应]
    H --> B
    G -->|验证通过| I[Service层处理]
    I --> J{检查用户名<br/>是否已存在}
    J -->|已存在| K[返回用户名<br/>已存在错误]
    K --> B
    J -->|不存在| L{检查邮箱<br/>是否已存在}
    L -->|已存在| M[返回邮箱<br/>已存在错误]
    M --> B
    L -->|不存在| N[bcrypt<br/>密码哈希处理]
    N --> O[创建User实体]
    O --> P[Repository<br/>写入数据库]
    P --> Q[返回成功响应]
    Q --> R([跳转到登录页面])`
    },
    {
        id: '5.1.2-登录功能',
        title: '登录功能实现流程',
        mermaid: `flowchart TD
    A([开始登录]) --> B[输入用户名和密码]
    B --> C[发送POST请求<br/>/api/v1/auth/login]
    C --> D[Controller接收请求<br/>绑定登录DTO]
    D --> E[Service层处理]
    E --> F{根据用户名<br/>查询用户}
    F -->|用户不存在| G[返回用户<br/>不存在错误]
    G --> B
    F -->|用户存在| H[bcrypt.Compare<br/>HashAndPassword]
    H --> I{密码是否正确}
    I -->|密码错误| J[返回密码错误]
    J --> B
    I -->|密码正确| K[golang-jwt<br/>生成JWT Token<br/>编码用户ID和用户名]
    K --> L[返回Token<br/>和用户基本信息]
    L --> M[前端存储Token<br/>Pinia + localStorage]
    M --> N[Axios请求拦截器<br/>自动携带Token]
    N --> O[路由守卫验证Token]
    O --> P([允许访问受保护页面])`
    },
    {
        id: '5.1.3-知识问答',
        title: '知识问答功能实现流程',
        mermaid: `flowchart TD
    A([开始问答]) --> B[输入问题内容]
    B --> C[发送POST请求<br/>/api/v1/ask<br/>问题内容+会话ID]
    C --> D[Service层接收问题]
    D --> E{会话ID是否为空?}
    E -->|为空| F[创建新AskSession<br/>标题:问题前20字符]
    E -->|不为空| G[使用现有会话]
    F --> H[保存用户问题<br/>AskMessage role=user]
    G --> H
    H --> I[知识检索逻辑<br/>基于关键词模糊匹配<br/>KnowledgePoint表查询]
    I --> J{找到相关<br/>知识点?}
    J -->|找到| K[计算置信度评分<br/>基于匹配数量和程度]
    J -->|未找到| L[设置置信度为0]
    K --> M[生成回答内容]
    L --> M
    M --> N[保存助手回答<br/>AskMessage role=assistant]
    N --> O[返回响应:<br/>回答内容+置信度<br/>+相关知识点列表]
    O --> P([前端展示完成])`
    },
    {
        id: '5.1.4-知识图谱浏览',
        title: '知识图谱浏览功能实现流程',
        mermaid: `flowchart TD
    A([进入图谱页面]) --> B[发送GET请求<br/>/api/v1/graph]
    B --> C[Controller处理请求]
    C --> D[查询KnowledgePoint表<br/>获取所有知识点]
    D --> E[查询KnowledgeRelation表<br/>获取所有关系]
    E --> F[组装JSON响应<br/>nodes + edges]
    F --> G[返回图谱数据]
    G --> H[前端接收数据]
    H --> I[D3.js forceSimulation<br/>创建力导向图模拟]
    I --> J[设置力:<br/>charge电荷力<br/>link链接力<br/>center中心力]
    J --> K[SVG渲染:<br/>知识点=圆形节点<br/>关系=连线+关系类型标注]
    K --> L([用户交互:<br/>拖拽节点+滚轮缩放])`
    },
    {
        id: '5.1.5-题库练习',
        title: '题库练习功能实现流程',
        mermaid: `flowchart TD
    A([进入题库页面]) --> B[发送GET请求<br/>/api/v1/question]
    B --> C[后端查询Question表<br/>返回题目列表]
    C --> D[前端展示题目卡片<br/>标题/类型/难度/知识点]
    D --> E[学生选择题目]
    E --> F[展示完整题目<br/>内容和选项]
    F --> G[学生选择答案]
    G --> H[发送POST请求<br/>/api/v1/quiz<br/>题目ID+用户答案]
    H --> I[Service层处理]
    I --> J[查询题目记录<br/>获取正确答案]
    J --> K[比对用户答案<br/>与正确答案]
    K --> L{答案是否一致?}
    L -->|一致| M[标记为正确<br/>is_correct = true]
    L -->|不一致| N[标记为错误<br/>is_correct = false]
    M --> O[保存答题记录<br/>到Quiz表]
    N --> O
    O --> P([返回结果:<br/>是否正确+正确答案<br/>+详细解析])`
    },
    {
        id: '5.1.6-学习分析',
        title: '学习分析功能实现流程',
        mermaid: `flowchart TD
    A([进入学习分析页面]) --> B[发送GET请求<br/>/api/v1/analytics/overview]
    B --> C[Service层聚合统计]
    C --> D[查询AskSession表<br/>统计问答会话总数]
    C --> E[查询AskMessage表<br/>统计消息总数]
    C --> F[查询Quiz表<br/>统计答题总数和正确数]
    C --> G[按难度统计正确率<br/>easy/medium/hard]
    C --> H[分析薄弱知识点<br/>基于答题记录]
    D --> I[封装统计数据]
    E --> I
    F --> I
    G --> I
    H --> I
    I --> J[返回统计数据]
    J --> K([ECharts渲染:<br/>问答统计+答题统计<br/>+薄弱知识点列表<br/>+综合雷达图])`
    },
    {
        id: '5.2.1-教师登录',
        title: '教师登录功能实现流程',
        mermaid: `flowchart TD
    A([教师开始登录]) --> B[输入教师账号和密码]
    B --> C[发送POST请求<br/>/api/v1/teacher_auth/login]
    C --> D[Controller接收请求<br/>绑定登录DTO]
    D --> E[Service层处理]
    E --> F[根据账号查询<br/>Teacher表]
    F --> G{教师是否存在?}
    G -->|不存在| H[返回教师<br/>不存在错误]
    H --> B
    G -->|存在| I[bcrypt<br/>比对密码]
    I --> J{密码是否正确?}
    J -->|错误| K[返回密码错误]
    K --> B
    J -->|正确| L[生成JWT Token<br/>包含教师身份标识]
    L --> M[RequireTeacherAuth<br/>中间件验证]
    M --> N[返回Token<br/>和教师信息]
    N --> O([进入教师管理界面])`
    },
    {
        id: '5.2.2-知识图谱构建',
        title: '知识图谱构建功能实现流程',
        mermaid: `flowchart TD
    A([教师点击构建图谱]) --> B[弹出资料选择对话框<br/>展示已审核文档列表]
    B --> C[勾选需要的文档]
    C --> D[点击确认构建]
    D --> E[发送POST请求<br/>/api/v1/graph/build<br/>文档ID列表]
    E --> F[Service层处理]
    F --> G[查询Document记录<br/>获取文档解析内容]
    G --> H[分块处理文档内容]
    H --> I[提取关键知识点<br/>和关系]
    I --> J[确定性规则抽取<br/>知识点实体和关系]
    J --> K[去重检查<br/>与KnowledgePoint表比对]
    K --> L{是否重复?}
    L -->|重复| M[跳过重复项]
    L -->|不重复| N[写入<br/>KnowledgePoint表]
    M --> O{还有更多<br/>知识点?}
    N --> P[写入<br/>KnowledgeRelation表]
    P --> O
    O -->|有| J
    O -->|没有| Q[记录构建统计信息<br/>到KnowledgeBuild表]
    Q --> R([返回构建结果:<br/>新增知识点数<br/>+新增关系数<br/>+跳过数])`
    }
];

// 生成HTML模板
function generateHTML(title, mermaidCode) {
    return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
    <style>
        body {
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
            padding: 20px;
            background: white;
            font-family: 'Microsoft YaHei', Arial, sans-serif;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 8px;
        }
        .title {
            text-align: center;
            font-size: 20px;
            font-weight: bold;
            margin-bottom: 20px;
            color: #000;
        }
        .mermaid {
            display: flex;
            justify-content: center;
        }
        .mermaid svg {
            max-width: 100%;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="title">${title}</div>
        <div class="mermaid">
            ${mermaidCode}
        </div>
    </div>
    <script>
        mermaid.initialize({
            startOnLoad: true,
            theme: 'base',
            themeVariables: {
                primaryColor: '#ffffff',
                primaryTextColor: '#000000',
                primaryBorderColor: '#000000',
                lineColor: '#000000',
                secondaryColor: '#f0f0f0',
                tertiaryColor: '#ffffff',
                fontFamily: 'Microsoft YaHei, Arial, sans-serif',
                fontSize: '14px'
            },
            flowchart: {
                useMaxWidth: false,
                htmlLabels: true,
                curve: 'basis'
            }
        });
    </script>
</body>
</html>`;
}

// 导出单个PNG
async function exportPNG(browser, flow, outputDir) {
    const page = await browser.newPage();
    const html = generateHTML(flow.title, flow.mermaid);

    await page.setViewport({ width: 1920, height: 1080 });
    await page.setContent(html, { waitUntil: 'networkidle0' });
    await page.waitForSelector('.mermaid svg', { timeout: 10000 });

    // 获取SVG元素的边界框
    const element = await page.$('.container');
    if (element) {
        const filePath = path.join(outputDir, `${flow.id}.png`);
        await element.screenshot({
            path: filePath,
            type: 'png',
            omitBackground: true
        });
        console.log(`✓ 已导出: ${flow.id}.png`);
    } else {
        console.error(`✗ 导出失败: ${flow.id}`);
    }

    await page.close();
}

// 主函数
async function main() {
    const outputDir = path.join(__dirname, 'flowcharts-png');

    // 创建输出目录
    if (!fs.existsSync(outputDir)) {
        fs.mkdirSync(outputDir, { recursive: true });
    }

    console.log('正在启动浏览器...');
    const browser = await puppeteer.launch({
        headless: 'new',
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });

    console.log(`开始导出 ${flows.length} 个流程图...\n`);

    for (const flow of flows) {
        await exportPNG(browser, flow, outputDir);
    }

    await browser.close();
    console.log(`\n导出完成！文件保存在: ${outputDir}`);
}

// 运行
main().catch(console.error);
